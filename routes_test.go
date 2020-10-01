package main

import (
	"bytes"
	"counter/store"
	"counter/store/mock_store"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	key := "/some/key"
	val := store.Value{
		Count: 5,
	}

	marshalled := marshal(key, &val)

	var result map[string]int

	err := json.Unmarshal([]byte(marshalled), &result)
	assert.NoError(t, err)

	assert.Equal(t, val.Count, result[key])
}

func TestRoutes_GetCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, uri, nil)

	rs := NewRoutes(repo)

	rs.GetCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, res.Body)
	assert.NoError(t, err)
	assert.Equal(t, marshal(uri, &v), buf.String())
}

func TestRoutes_GetCounter404(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     0,
		AccessKey: uuid.Nil,
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, uri, nil)

	rs := NewRoutes(repo)

	rs.GetCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestRoutes_authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, uri, nil)
	r.Header.Set("Authorization", "Bearer "+v.AccessKey.String())

	rs := NewRoutes(repo)

	assert.True(t, rs.authenticate(w, r))
}

func TestRoutes_authenticate_noheader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, uri, nil)

	rs := NewRoutes(repo)

	assert.False(t, rs.authenticate(w, r))
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestRoutes_authenticate_wrongheader(t *testing.T) {
	// Bearer but invalid uuid
	helperAuthenticateWrongheader(t, "Bearer Yeet")
	// No bearer
	helperAuthenticateWrongheader(t, "Yeet")
	// Wrong UUID
	helperAuthenticateWrongheader(t, "Bearer "+uuid.New().String())
}

func helperAuthenticateWrongheader(t *testing.T, header string) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, uri, nil)
	r.Header.Set("Authorization", header)

	rs := NewRoutes(repo)

	assert.False(t, rs.authenticate(w, r))
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestRoutes_PatchCounter(t *testing.T){
	PatchTests(t,"increment")
	PatchTests(t,"decrement")
}

func PatchTests(t *testing.T, op string) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)
	if op == "increment" {
		repo.EXPECT().Increment(uri).Return(nil).Times(1)
	} else if op == "decrement" {
		repo.EXPECT().Decrement(uri).Return(nil).Times(1)
	} else {
		t.Fail()
	}

	w := httptest.NewRecorder()
	b, err := json.Marshal(patchArgs{Op: op})
	assert.NoError(t, err)
	r := httptest.NewRequest(http.MethodPatch, uri, bytes.NewReader(b))
	r.Header.Set("Authorization", "Bearer "+v.AccessKey.String())

	rs := NewRoutes(repo)

	rs.PatchCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)
	assert.NoError(t, err)
	assert.Equal(t, marshal(uri, &v), buf.String())
}

func TestRoutes_IncrementCounterUnAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, uri, nil)

	rs := NewRoutes(repo)

	rs.PatchCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestRoutes_IncrementCounterNX(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.Nil,
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, uri, nil)

	rs := NewRoutes(repo)

	rs.PatchCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestRoutes_CreateCounter_Exists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, uri, nil)
	r.Header.Set("Authorization", "Bearer "+v.AccessKey.String())

	rs := NewRoutes(repo)

	rs.CreateCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusConflict, res.StatusCode)
}

func TestRoutes_CreateCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     0,
		AccessKey: uuid.Nil,
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).Times(1)
	repo.EXPECT().Create(uri, gomock.Any()).Return(nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, uri, nil)

	rs := NewRoutes(repo)

	rs.CreateCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	header := w.Header().Get("Authorization")
	assert.True(t, strings.HasPrefix(header, "Bearer "))
	token := strings.TrimPrefix(header, "Bearer ")
	id, err := uuid.Parse(token)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id.String())

	buf := new(strings.Builder)
	_, err = io.Copy(buf, res.Body)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), id.String())
	assert.Contains(t, buf.String(), uri)
}

func TestRoutes_DeleteCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     0,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)
	repo.EXPECT().Delete(uri).Return(nil).Times(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, uri, nil)
	r.Header.Set("Authorization", "Bearer "+v.AccessKey.String())

	rs := NewRoutes(repo)

	rs.DeleteCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRoutes_DeleteCounterUnAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     0,
		AccessKey: uuid.New(),
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, uri, nil)

	rs := NewRoutes(repo)

	rs.DeleteCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestRoutes_DeleteCounterNX(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	v := store.Value{
		Count:     0,
		AccessKey: uuid.Nil,
	}

	uri := "/yeet"

	repo := mock_store.NewMockRepository(ctrl)

	repo.EXPECT().Get(uri).Return(v, nil).MinTimes(1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, uri, nil)
	r.Header.Set("Authorization", "Bearer "+v.AccessKey.String())

	rs := NewRoutes(repo)

	rs.DeleteCounter(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
