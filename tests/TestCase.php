<?php

namespace App\Tests;

use PHPUnit\Framework\MockObject\MockObject;

class TestCase extends \PHPUnit\Framework\TestCase
{
    public function getMock(string $className, array $methods = []): MockObject
    {
        return $this->getMockBuilder($className)
                    ->disableArgumentCloning()
                    ->disableProxyingToOriginalMethods()
                    ->disableOriginalConstructor()
                    ->setMethods($methods)
                    ->getMock();
    }
}