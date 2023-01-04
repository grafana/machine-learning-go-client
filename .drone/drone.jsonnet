local step(name, commands, image='golang:1.17') = {
  name: name,
  commands: commands,
  image: image,
};

local pipeline(name, steps=[]) = {
  kind: 'pipeline',
  type: 'docker',
  name: name,
  steps: steps,
  trigger+: {
    ref+: [
      'refs/heads/main',
      'refs/pull/**',
      'refs/tags/**',
    ],
  },
};

[
  pipeline('Lint', steps=[
    step('lint',
         ['golangci-lint run'],
         'golangci/golangci-lint:v1.50.1'),
  ]),
  pipeline('Test', steps=[
    step('test', ['go test -cover -v ./...']),
  ]),
]
