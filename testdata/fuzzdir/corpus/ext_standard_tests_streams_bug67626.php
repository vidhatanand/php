<?php
class MyStream
{
	public function stream_open() { return true; }
	
	public function stream_read()
	{
		throw new Exception('stream_read_exception');
		return 'read';
	}
	
	public function stream_eof()
	{
		return true;
	}
	
	public function stream_write()
	{
		throw new Exception('stream_write_exception');
		return 42;
	}
}

stream_wrapper_register("my", "MyStream");

$fp = fopen('my://foobar', 'r+');

try {
	fread($fp, 42);
} catch (Exception $e) {
	echo $e->getMessage();
}
echo "\n";
try {
	fwrite($fp, 'foobar');
} catch (Exception $e) {
	echo $e->getMessage();
}
?>
