import { defineConfig } from "drizzle-kit";

const dbUrl = `postgresql://${process.env.DB_USER}:${process.env.DB_PASSWORD}@${process.env.DB_HOST}:${process.env.DB_PORT}/${process.env.DB_NAME}?schema=public`;

export default defineConfig({
	schema: "./migrations/schema.ts",
	out: "./migrations/drizzle",
	dialect: "postgresql",
	dbCredentials: {
		url: dbUrl,
	},
	introspect: {
		casing: "preserve",
	},
});
