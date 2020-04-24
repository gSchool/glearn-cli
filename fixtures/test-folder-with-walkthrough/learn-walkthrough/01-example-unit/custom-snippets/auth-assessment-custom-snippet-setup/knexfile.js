console.log(`From knexfile`)
console.log(process.env.DATABASE_URL)
console.log(__dirname)

module.exports = {
  development: {
    client: 'pg',
    connection: 'postgres://localhost:5432/auth',
    migrations: {
      directory: __dirname + '/src/db/migrations'
    },
    seeds: {
      directory: __dirname + '/src/db/seeds'
    }
  },
  test: {
    client: 'pg',
    connection: process.env.DATABASE_URL,
    migrations: {
      directory: __dirname + '/src/db/migrations'
    },
    seeds: {
      directory: __dirname + '/src/db/seeds'
    },
    pool: {
      min: 0,
      max: 1
    }
  }
};
