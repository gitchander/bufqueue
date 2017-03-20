package bufqueue

// Block structure:
// +--------+--------+--------+---------+--------+
// | record | record | record |   ...   | record |
// +--------+--------+--------+---------+--------+

// Record structure:
// +----------+-------------+
// |  Length  |   Message   |
// +----------+-------------+

type FilesConfig struct {
	Dirname    string
	FilePrefix string
}
