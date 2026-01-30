import {
	bigint,
	boolean,
	date,
	integer,
	json,
	pgTable,
	text,
	time,
	timestamp,
	varchar,
} from "drizzle-orm/pg-core";

export const soccerMatches = pgTable("soccer_matches", {
	id: bigint("id", { mode: "number" }).primaryKey().generatedAlwaysAsIdentity(),
	matchId: bigint("match_id", { mode: "number" }).unique(),
	leagueGid: bigint("league_gid", { mode: "number" }),
	leagueId: bigint("league_id", { mode: "number" }),
	leagueName: varchar("league_name", { length: 255 }),
	matchStatus: varchar("match_status", { length: 50 }),
	matchStartDate: date("match_start_date"),
	matchStartTime: time("match_start_time"),
	hTeamId: bigint("h_team_id", { mode: "number" }),
	aTeamId: bigint("a_team_id", { mode: "number" }),
	hTeamName: varchar("h_team_name", { length: 255 }),
	aTeamName: varchar("a_team_name", { length: 255 }),
	hTeamGoals: integer("h_team_goals"),
	aTeamGoals: integer("a_team_goals"),
	htScore: varchar("ht_score", { length: 10 }),
	ftScore: varchar("ft_score", { length: 10 }),
	events: json("events"),
	createdAt: timestamp("created_at").defaultNow(),
	updatedAt: timestamp("updated_at").defaultNow(),
});

export const basketballMatches = pgTable("basketball_matches", {
	id: bigint("id", { mode: "number" }).primaryKey().generatedAlwaysAsIdentity(),
	matchId: bigint("match_id", { mode: "number" }).unique(),
	leagueGid: bigint("league_gid", { mode: "number" }),
	leagueId: bigint("league_id", { mode: "number" }),
	leagueName: varchar("league_name", { length: 255 }),
	fileGroup: varchar("file_group", { length: 100 }),
	matchStatus: varchar("match_status", { length: 50 }),
	matchDate: date("match_date"),
	matchTime: time("match_time"),
	timer: varchar("timer", { length: 20 }),
	hTeamId: bigint("h_team_id", { mode: "number" }),
	hTeamName: varchar("h_team_name", { length: 255 }),
	hTeamScore: integer("h_team_score"),
	hTeamQ1: integer("h_team_q1"),
	hTeamQ2: integer("h_team_q2"),
	hTeamQ3: integer("h_team_q3"),
	hTeamQ4: integer("h_team_q4"),
	hTeamOt: integer("h_team_ot"),
	aTeamId: bigint("a_team_id", { mode: "number" }),
	aTeamName: varchar("a_team_name", { length: 255 }),
	aTeamScore: integer("a_team_score"),
	aTeamQ1: integer("a_team_q1"),
	aTeamQ2: integer("a_team_q2"),
	aTeamQ3: integer("a_team_q3"),
	aTeamQ4: integer("a_team_q4"),
	aTeamOt: integer("a_team_ot"),
	createdAt: timestamp("created_at").defaultNow(),
	updatedAt: timestamp("updated_at").defaultNow(),
});

// API Keys for authentication
export const apiKeys = pgTable("api_keys", {
	id: bigint("id", { mode: "number" }).primaryKey().generatedAlwaysAsIdentity(),
	keyHash: varchar("key_hash", { length: 64 }).notNull().unique(), // SHA256 hash
	keyPrefix: varchar("key_prefix", { length: 16 }).notNull(), // "sk_live_xxxx" for display
	name: varchar("name", { length: 255 }).notNull(), // Human-readable name
	sports: json("sports").notNull(), // Array of allowed sports: ["soccer", "basketball"] or ["*"]
	rateLimit: integer("rate_limit").notNull().default(100), // Requests per minute
	isActive: boolean("is_active").notNull().default(true),
	createdAt: timestamp("created_at").defaultNow(),
	lastUsedAt: timestamp("last_used_at"),
	expiresAt: timestamp("expires_at"),
});
