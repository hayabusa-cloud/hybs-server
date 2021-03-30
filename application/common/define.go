package common

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
)

const (
	UserPermissionGuest  = hybs.UserPermissionGuest
	UserPermissionNormal = hybs.UserPermissionNormal
	UserPermissionStaff  = UserPermissionNormal + 1
	UserPermissionAdmin  = UserPermissionStaff + 1
)

const (
	TimeLayoutDefault = "2006-01-02 15:04"
	TimeLayoutYMDHm   = "2006-01-02 15:04"
	TimeLayoutYMDHms  = "2006-01-02 15:04:05"
)

const (
	FloatEps = 1e-6
)
