package tracing

import (
	"github.com/getsentry/sentry-go"
	"gorm.io/gorm"
)

type SentryPlugin struct{}

func NewSentryPlugin() *SentryPlugin {
	return &SentryPlugin{}
}

func (p *SentryPlugin) Name() string {
	return "sentry-plugin"
}

func (p *SentryPlugin) Initialize(db *gorm.DB) error {
	// Register callbacks for errors and queries
	db.Callback().Query().Before("gorm:query").Register("sentry:before_query", p.startSpan)
	db.Callback().Query().After("gorm:query").Register("sentry:after_query", p.endSpan)

	db.Callback().Row().Before("gorm:row").Register("sentry:before_row", p.startSpan)
	db.Callback().Row().After("gorm:row").Register("sentry:after_row", p.endSpan)

	db.Callback().Create().Before("gorm:create").Register("sentry:before_create", p.startSpan)
	db.Callback().Create().After("gorm:create").Register("sentry:after_create", p.endSpan)

	db.Callback().Update().Before("gorm:update").Register("sentry:before_update", p.startSpan)
	db.Callback().Update().After("gorm:update").Register("sentry:after_update", p.endSpan)

	db.Callback().Delete().Before("gorm:delete").Register("sentry:before_delete", p.startSpan)
	db.Callback().Delete().After("gorm:delete").Register("sentry:after_delete", p.endSpan)
	return nil
}

func (p *SentryPlugin) startSpan(db *gorm.DB) {
	transaction, ok := db.Statement.Context.Value("sentry_transaction").(*sentry.Span)
	if !ok {
		return
	}

	span := sentry.StartSpan(transaction.Context(), "db.query")
	span.Data = map[string]interface{}{
		"sql":          db.Statement.SQL.String(),
		"table":        db.Statement.Table,
		"rowsAffected": db.Statement.RowsAffected,
	}
	db.InstanceSet("sentry:span", span)
}

func (p *SentryPlugin) endSpan(db *gorm.DB) {
	if v, ok := db.InstanceGet("sentry:span"); ok {
		if span, ok := v.(*sentry.Span); ok {
			span.Finish()
		}
	}
}
