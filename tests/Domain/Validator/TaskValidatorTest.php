<?php

namespace App\Tests\Domain\Validator;

use App\Domain\Dto\TaskDto;
use App\Domain\Exception\InvalidTaskException;
use App\Domain\Validator\TaskValidator;
use App\Tests\TestCase;

class TaskValidatorTest extends TestCase
{

    /**
     * @var TaskValidator
     */
    private $validator;

    protected function setUp()
    {
        parent::setUp();
        $this->validator = new TaskValidator();
    }

    public function testValidateAllValid()
    {
        $taskDto = new TaskDto();
        $taskDto->name = "name";
        $taskDto->description = "description";
        $taskDto->when = (new \DateTime())->modify("+1 day")->format('Y-m-d');
        $this->assertNull($this->validator->validate($taskDto));
    }

    public function testValidateInvalidDataAndNoDate()
    {
        $taskDto = new TaskDto();
        $this->expectExceptionObject(new InvalidTaskException([]));
        $this->validator->validate($taskDto);
    }

    public function testValidateInvalidDate()
    {
        $taskDto = new TaskDto();
        $taskDto->name = "name";
        $taskDto->description = "description";
        $taskDto->when = (new \DateTime())->modify("+1 day")->format('Y-m');
        $this->expectExceptionObject(new InvalidTaskException([]));
        $this->validator->validate($taskDto);
    }

    public function testValidatePastDate()
    {
        $taskDto = new TaskDto();
        $taskDto->name = "name";
        $taskDto->description = "description";
        $taskDto->when = (new \DateTime())->modify("-1 day")->format('Y-m-d');
        $this->expectExceptionObject(new InvalidTaskException([]));
        $this->validator->validate($taskDto);
    }
}