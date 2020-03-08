<?php

namespace App\Tests\Domain\Service;

use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Repository\TaskRepositoryInterface;
use App\Domain\Service\TaskService;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Domain\Validator\TaskValidator;
use App\Tests\TestCase;

class TaskServiceTest extends TestCase
{
    private $repositoryMock;

    private $validatorMock;

    private $transformerMock;

    protected function setUp()
    {
        parent::setUp();
        $this->repositoryMock = $this->getMock(TaskRepositoryInterface::class, ['saveTask', 'findTaskById', 'findAll']);
        $this->validatorMock = $this->getMock(TaskValidator::class, ['validate']);
        $this->transformerMock = $this->getMock(TaskDtoToEntityTransformer::class, ['transform', 'reverseTransform']);
    }

    public function testAddTask()
    {
        $taskMock = $this->getMock(Task::class);
        $taskDto = new TaskDto();
        $this->repositoryMock->expects(self::once())->method('saveTask')->with($taskMock)->willReturn($taskMock);
        $this->validatorMock->expects(self::once())->method('validate')->with($taskDto)->willReturn(true);
        $this->transformerMock->expects(self::once())->method('transform')->with($taskDto)->willReturn($taskMock);
        $this->transformerMock->expects(self::once())->method('reverseTransform')->with($taskMock)->willReturn(
            $taskDto
        );
        $service = new TaskService($this->repositoryMock, $this->validatorMock, $this->transformerMock);
        $this->assertSame($taskDto, $service->addTask($taskDto));
    }
}
