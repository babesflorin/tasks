<?php

namespace App\Presentation\EventListener;

use App\Domain\Exception\ValidationException;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\ExceptionEvent;

class ExceptionListener
{
    public function onKernelException(ExceptionEvent $event)
    {
        // You get the exception object from the received event
        $exception = $event->getThrowable();

        $statusCode = Response::HTTP_INTERNAL_SERVER_ERROR;
        // Customize your response object to display the exception details
        $responseArray = ['data' => '', 'error' => $exception->getMessage()];
        if ($exception instanceof ValidationException) {
            $responseArray['messages'] = $exception->getErrors();
            $statusCode = Response::HTTP_BAD_REQUEST;
        }
        $response = new JsonResponse($responseArray, $statusCode);

        // sends the modified response object to the event
        $event->setResponse($response);
    }
}
