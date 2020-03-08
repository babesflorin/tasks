<?php

namespace App\Domain\Service;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Exception\InvalidTaskException;
use App\Domain\Repository\TaskRepositoryInterface;
use App\Domain\Transformer\TaskCollectionDtoToCollectionTransformer;
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
    private $dtoTransformer;
    /**
     * @var TaskCollectionDtoToCollectionTransformer
     */
    private $dtoCollectionTransformer;

    public function __construct(
        TaskRepositoryInterface $repository,
        TaskValidator $validator,
        TaskDtoToEntityTransformer $dtoTransformer,
        TaskCollectionDtoToCollectionTransformer $dtoCollectionTransformer
    ) {
        $this->repository = $repository;
        $this->validator = $validator;
        $this->dtoTransformer = $dtoTransformer;
        $this->dtoCollectionTransformer = $dtoCollectionTransformer;
    }

    /**
     * @throws InvalidTaskException
     */
    public function addTask(TaskDto $taskDto): TaskDto
    {
        $this->validator->validate($taskDto);
        $task = $this->repository->saveTask($this->dtoTransformer->transform($taskDto));

        return $this->dtoTransformer->reverseTransform($task);
    }

    public function getAllTasks(): TaskCollectionDto
    {
        $tasks = $this->repository->getTasks();

        return $this->dtoCollectionTransformer->reverseTransform($tasks);
    }
}