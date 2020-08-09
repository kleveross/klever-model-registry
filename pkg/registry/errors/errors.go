package errors

import (
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/caicloud/nirvana/errors"
)

// RenderError is convert k8s error to nirvana error format, and it is for http response.
func RenderError(err error) error {
	if k8serrors.IsAlreadyExists(err) || k8serrors.IsInvalid(err) {
		return RenderBadRequestError(err)
	} else if k8serrors.IsNotFound(err) {
		return RenderNotFoundError(err)
	} else if k8serrors.IsUnauthorized(err) {
		return RenderUnAuthorizedError(err)
	} else if k8serrors.IsConflict(err) {
		return RenderSendConflictError(err)
	} else if k8serrors.IsInternalError(err) {
		return RenderInternalServerError(err)
	} else if k8serrors.IsForbidden(err) {
		return RenderForbiddenError(err)
	} else if k8serrors.IsNotAcceptable(err) {
		return RenderNotAcceptableError(err)
	} else if k8serrors.IsServerTimeout(err) {
		return RenderRequestTimeoutError(err)
	}

	return RenderInternalServerError(err)
}

func RenderBadRequestError(err error) error {
	return errors.BadRequest.Error(err.Error())
}

func RenderNotFoundError(err error) error {
	return errors.NotFound.Error(err.Error())
}

func RenderUnAuthorizedError(err error) error {
	return errors.Unauthorized.Error(err.Error())
}

func RenderSendConflictError(err error) error {
	return errors.Conflict.Error(err.Error())
}

func RenderInternalServerError(err error) error {
	return errors.InternalServerError.Error(err.Error())
}

func RenderForbiddenError(err error) error {
	return errors.Forbidden.Error(err.Error())
}

func RenderPreconditionFailedError(err error) error {
	return errors.PreconditionFailed.Error(err.Error())
}

func RenderStatusServiceUnavailableError(err error) error {
	return errors.ServiceUnavailable.Error(err.Error())
}

func RenderNotAcceptableError(err error) error {
	return errors.NotAcceptable.Error(err.Error())
}

func RenderRequestTimeoutError(err error) error {
	return errors.RequestTimeout.Error(err.Error())
}
