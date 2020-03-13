<?php

namespace App\Tests\Unit\Domain\Entity;

use App\Domain\Entity\Task;
use App\Tests\Unit\TestCase;

class TaskTest extends TestCase
{

    public function testComplete()
    {
        $task = new Task("name", "description", new \DateTime());
        $this->assertFalse($task->isDone());
        $task->complete();
        $this->assertTrue($task->isDone());
    }
}
