package repository

import "os"

const ID_NOT_FOUND = "ID specified was not found"

// this one can't be a constant, unless we agree on a static PORT number for dev
var BASE_URL = "http://localhost:" + os.Getenv("PORT") + "/"
