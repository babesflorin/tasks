<?php

namespace App\Domain\Service;

use App\Domain\Dto\TaskDto;
use App\Domain\Exception\InvalidTaskException;
use App\Domain\Repository\TaskRepositoryInterface;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Domain\Validator\TaskValidator;

class TaskService
{
    /**
     * @var TaskRepositoryInterface
     */
    private $repository;
    /**
     * @var TaskValidator
     */
    private $validator;
    /**
     * @var TaskDtoToEntityTransformer
     */
    private $dtoToEntityTransfomer;

    public function __construct(
        TaskRepositoryInterface $repository,
        TaskValidator $validator,
        TaskDtoToEntityTransformer $dtoToEntityTransfomer
    ) {
        $this->repository = $repository;
        $this->validator = $validator;
        $this->dtoToEntityTransfomer = $dtoToEntityTransfomer;
    }

    /**
     * @throws InvalidTaskException
     */
    public function addTask(TaskDto $taskDto): TaskDto
    {
        $this->validator->validate($taskDto);
        $task = $this->repository->saveTask($this->dtoToEntityTransfomer->transform($taskDto));

        return $this->dtoToEntityTransfomer->reverseTransform($task);
    }
}