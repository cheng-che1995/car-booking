package main

import (
	"testing"
)

func TestValidate(t *testing.T) {
	user1 := &User{Username: "", Password: ""}
	user2 := &User{Username: "", Password: "123"}
	user3 := &User{Username: "user3", Password: ""}
	user4 := &User{Username: "user4", Password: "123455"}

	if case1 := user1.Validate(); case1 != ErrUsernameEmpty {
		t.Errorf("Expected 使用者名稱不得為空值！ but got %v", case1)
	}
	if case2 := user2.Validate(); case2 != ErrUsernameEmpty {
		t.Errorf("Expected 使用者名稱不得為空值！ but got %v", case2)
	}
	if case3 := user3.Validate(); case3 != ErrPasswordEmpty {
		t.Errorf("Expected 密碼不得為空值！ but got %v", case3)
	}
	if case4 := user4.Validate(); case4 != nil {
		t.Errorf("Expected nil but got %v", case4)
	}
}
