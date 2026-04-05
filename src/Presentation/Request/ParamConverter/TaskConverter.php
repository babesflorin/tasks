<?php

namespace App\Presentation\Request\ParamConverter;

use App\Domain\Dto\TaskDto;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\ParamConverter;
use Sensio\Bundle\FrameworkExtraBundle\Request\ParamConverter\ParamConverterInterface;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Exception\HttpException;

class TaskConverter implements ParamConverterInterface
{

    /**
     * @inheritDoc
     */
    public function apply(Request $request, ParamConverter $configuration)
    {
        $data = json_decode($request->getContent(), true);
        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new HttpException(Response::HTTP_BAD_REQUEST, "Request must be json!");
        }
        $task = new TaskDto();
        $task->id = $data['id'] ?? null;
        $task->name = $data['name'] ?? '';
        $task->description = $data['description'] ?? '';
        $task->when = $data['when'] ?? '';
        $request->attributes->set($configuration->getName(), $task);
    }

    /**
     * @inheritDoc
     */
    public function supports(ParamConverter $configuration)
    {
        return $configuration->getClass() === TaskDto::class;
    }
}
