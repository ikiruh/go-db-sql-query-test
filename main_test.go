package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	conn, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer conn.Close()

	clientID := 1
	client, err := selectClient(conn, clientID)
	require.NoError(t, err)
	assert.Equal(t, clientID, client.ID)
	assert.NotEmpty(t, client.FIO)
	assert.NotEmpty(t, client.Login)
	assert.NotEmpty(t, client.Birthday)
	assert.NotEmpty(t, client.Email)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	conn, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer conn.Close()

	clientID := -1
	client, err := selectClient(conn, clientID)
	require.ErrorIs(t, err, sql.ErrNoRows)

	assert.Empty(t, client.ID)
	assert.Empty(t, client.Login)
	assert.Empty(t, client.FIO)
	assert.Empty(t, client.Birthday)
	assert.Empty(t, client.Email)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	conn, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer conn.Close()

	insClient := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	insClient.ID, err = insertClient(conn, insClient)
	require.NoError(t, err)
	require.NotEmpty(t, insClient.ID)

	selClient, err := selectClient(conn, insClient.ID)
	require.NoError(t, err)
	assert.Equal(t, insClient, selClient)
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	conn, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer conn.Close()

	insClient := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	insClient.ID, err = insertClient(conn, insClient)
	require.NoError(t, err)
	require.NotEmpty(t, insClient.ID)

	_, err = selectClient(conn, insClient.ID)
	require.NoError(t, err)

	err = deleteClient(conn, insClient.ID)
	require.NoError(t, err)

	_, err = selectClient(conn, insClient.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
