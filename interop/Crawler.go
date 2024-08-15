//go:build interop
package Interop

/*
#cgo CFLAGS: -I/usr/include/python3.11
#cgo LDFLAGS: -lpython3.11

#include <Python.h> // Python C-API
#include <stdlib.h>
#include <stdio.h>

char* load_cpython()
{
	if (Py_IsInitialized())
		return NULL;

	Py_Initialize();
	if (!Py_IsInitialized())
		return "Failed to initialize Python. Py_Initialize did not set up Python interpreter.";

	PyObject* sysPath = PySys_GetObject("path");
	if (sysPath == NULL)
		return "Failed to get sys.path";

	PyObject* cwd = PyUnicode_DecodeFSDefault(".");
	if (cwd == NULL)
		return "Failed to decode current directory";

	if (PyList_Append(sysPath, cwd) != 0)
		return "Failed to append current directory to sys.path";

	Py_DECREF(cwd);
	return NULL;
}

char* handle_error()
{
	PyObject *type, *value, *traceback;
	PyErr_Fetch(&type, &value, &traceback);
	const char* error_name = PyExceptionClass_Name(type);

	PyObject* value_of_error_obj = PyObject_Str(value);
	PyObject* bytes_utf8_value = PyUnicode_AsUTF8String(value_of_error_obj);
	char* value_as_c_string = PyBytes_AsString(bytes_utf8_value);

	PyObject* traceback_str = PyObject_Str(traceback);
	PyObject* bytes_utf8_traceback = PyUnicode_AsUTF8String(traceback_str);
	char* traceback_as_c_string = PyBytes_AsString(bytes_utf8_traceback);

	char* res = malloc(strlen(error_name) + strlen(value_as_c_string) + strlen(traceback_as_c_string) + 4);
	sprintf(res, "%s: %s\n%s", error_name, value_as_c_string, traceback_as_c_string);

	Py_DECREF(type);
	Py_DECREF(value);
	Py_DECREF(traceback);
	Py_DECREF(value_of_error_obj);
	Py_DECREF(traceback_str);

	PyErr_Clear();
	return res;
}

// Get links from URL
char** extract_links_from_url(char* url, int depth, char** out_error)
{
	PyObject* pName = NULL;
	PyObject* pModule = NULL;
	PyObject* pFunc = NULL;
	PyObject* pArgs = NULL;
	PyObject* pValue = NULL;
	char** result = NULL;

	PyGILState_STATE gstate;
	gstate = PyGILState_Ensure();

	pName = PyUnicode_DecodeFSDefault("crawler");
	pModule = PyImport_Import(pName);
	Py_DECREF(pName);

	if (pModule == NULL) // handle error
	{
		*out_error = PyErr_Occurred() ? handle_error() : strdup("Failed to load module crawler");
		PyErr_Clear();
		goto cleanup;
	}

	pFunc = PyObject_GetAttrString(pModule, "extract_links_from_url");
	if (!pFunc || !PyCallable_Check(pFunc)) // handle error
	{
		*out_error = PyErr_Occurred() ? handle_error() : strdup("Cannot find function extract_links_from_url");
		PyErr_Clear();
		goto cleanup;
	}

	pArgs = PyTuple_New(2);
	PyTuple_SetItem(pArgs, 0, PyUnicode_FromString(url)); // url parameter
	PyTuple_SetItem(pArgs, 1, PyLong_FromLong(depth)); // depth parameter

	pValue = PyObject_CallObject(pFunc, pArgs);
	Py_DECREF(pArgs);

	if (PyErr_Occurred() || pValue == NULL)
	{
		*out_error = PyErr_Occurred() ? handle_error() : strdup("function extract_links_from_url failed");
		PyErr_Clear();
		goto cleanup;
	}

	if (!PyList_Check(pValue))
	{
		*out_error = strdup("function extract_links_from_url did not return a list");
		goto cleanup;
	}

	Py_ssize_t size = PyList_Size(pValue); // get the size of the list
	result = malloc((size + 1) * sizeof(char*));
	result[size] = NULL; // mark the last element using NULL

	// copy the strings from the list to the C array
	for (Py_ssize_t i = 0; i < size; i++)
	{
		PyObject *item = PyList_GetItem(pValue, i); // the i-th string
		if (!PyUnicode_Check(item))
		{
			*out_error = strdup("function extract_links_from_url returned a non-string item");
			for (Py_ssize_t j = 0; j < i; j++)
			{
				free(result[j]);
			}
			free(result);
			goto cleanup;
		}
		PyObject* item_as_utf8 = PyUnicode_AsUTF8String(item); // convert to bytes as utf-8
		result[i] = strdup(PyBytes_AsString(item_as_utf8)); // copy the bytes to a new string
		Py_DECREF(item_as_utf8); // free the bytes object
	}

cleanup:
	Py_XDECREF(pFunc);
	Py_XDECREF(pModule);
	Py_XDECREF(pName);
	Py_XDECREF(pArgs);
	Py_XDECREF(pValue);
	PyGILState_Release(gstate);
	return result;
}

char* get_element(char** array, int index)
{
	return array[index];
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

func LoadPython() error {
	err := C.load_cpython()
	if err != nil {
		return errors.New(C.GoString(err))
	}
	return nil
}

func ExtractLinksFromURL(url string, depth int) ([]string, error) {
	c_url := C.CString(url)
	defer C.free(unsafe.Pointer(c_url))
	var c_error *C.char

	// call the C function
	c_result := C.extract_links_from_url(c_url, C.int(depth), &c_error)
	if c_error != nil {
		defer C.free(unsafe.Pointer(c_error))
		return nil, errors.New(C.GoString(c_error))
	}
	defer C.free(unsafe.Pointer(c_result))

	length := 0
	for C.get_element(c_result, C.int(length)) != nil {
		length++
	}

	// create a slice
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(c_result))[:length:length]

	// create the []string that will hold the result
	goStrings := make([]string, length)

	// copy the strings from the C array to the Go slice
	for i, s := range tmpslice {
		goStrings[i] = C.GoString(s)
		C.free(unsafe.Pointer(s))
	}

	return goStrings, nil
}
