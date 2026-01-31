package analysis

// Topic represents predefined business topics for categorizing feedback.
// These topics are used to classify feedback into specific business categories
// rather than allowing the LLM to create topics from scratch.
type Topic string

const (
	// TopicProductFunctionalityFeatures covers missing features, feature requests,
	// feature quality, incorrect behavior, logic bugs, and edge cases.
	TopicProductFunctionalityFeatures Topic = "product_functionality_features"

	// TopicUIUX covers layout, visual design, accessibility, navigation clarity,
	// learnability, consistency, and mobile vs desktop experience.
	TopicUIUX Topic = "ui_ux"

	// TopicPerformanceReliability covers app speed, memory usage, CPU/battery consumption,
	// network efficiency, crashes, freezes, uptime, and scalability.
	TopicPerformanceReliability Topic = "performance_reliability"

	// TopicUsabilityProductivity covers workflow efficiency, automation needs,
	// keyboard shortcuts, defaults, smart suggestions, and cognitive load.
	TopicUsabilityProductivity Topic = "usability_productivity"

	// TopicSecurityPrivacy covers data leaks, privacy concerns, authentication,
	// authorization, permission models, encryption, compliance, and token handling.
	TopicSecurityPrivacy Topic = "security_privacy"

	// TopicCompatibilityIntegration covers OS compatibility, browser compatibility,
	// API integrations, import/export formats, backward compatibility, and SDK quality.
	TopicCompatibilityIntegration Topic = "compatibility_integration"

	// TopicDeveloperExperience covers documentation quality, API ergonomics,
	// error messages clarity, debugging tools, SDK stability, code examples,
	// tutorials, and language support.
	TopicDeveloperExperience Topic = "developer_experience"

	// TopicPricingLicensing covers pricing concerns, confusing pricing tiers,
	// hidden costs, licensing restrictions, value-for-money perception, and free tier limitations.
	TopicPricingLicensing Topic = "pricing_licensing"

	// TopicCustomerSupportCommunity covers support response time, quality of support answers,
	// community activity, bug triage speed, and transparency in roadmap.
	TopicCustomerSupportCommunity Topic = "customer_support_community"

	// TopicInstallationSetupDeployment covers installation difficulties, broken installers,
	// dependency issues, Docker images, CI/CD integration problems, and upgrade/migration issues.
	TopicInstallationSetupDeployment Topic = "installation_setup_deployment"

	// TopicDataAnalyticsReporting covers incorrect metrics, missing dashboards,
	// hard-to-export reports, real-time vs delayed data, and visualization issues.
	TopicDataAnalyticsReporting Topic = "data_analytics_reporting"

	// TopicLocalizationInternationalization covers poor translations, missing languages,
	// date/time/number formatting issues, RTL layout bugs, and cultural UX mismatches.
	TopicLocalizationInternationalization Topic = "localization_internationalization"

	// TopicProductStrategyRoadmap covers direction dissatisfaction, feature prioritization complaints,
	// requests for enterprise vs consumer features, and backward compatibility decisions.
	TopicProductStrategyRoadmap Topic = "product_strategy_roadmap"
)

// String returns the string representation of the topic.
func (t Topic) String() string {
	return string(t)
}

// Description returns the detailed description of the topic for use in LLM prompts.
func (t Topic) Description() string {
	switch t {
	case TopicProductFunctionalityFeatures:
		return `Product Functionality & Features

This is the most obvious category.

Topics here:
- Missing features ("I need bulk export")
- Feature requests ("Add dark mode", "Add GitHub integration")
- Feature quality ("Search is inaccurate", "Filters are confusing")
- Incorrect behavior / logic bugs
- Edge cases not handled

👉 This is **what the product does**.`

	case TopicUIUX:
		return `UI / UX (Interaction & Design)

This is broad and covers multiple aspects.

Subtopics:
- Layout & visual design
- Accessibility (contrast, screen readers, keyboard navigation)
- Navigation clarity (too many clicks, confusing menus)
- Learnability (onboarding, discoverability)
- Consistency (buttons behave differently in different screens)
- Mobile vs desktop experience

👉 This is **how the product feels to use**.`

	case TopicPerformanceReliability:
		return `Performance & Reliability

This is bigger than just performance.

Topics:
- App speed (startup, response times)
- Memory usage
- CPU/battery consumption
- Network efficiency (too many API calls)
- Crashes and freezes
- Uptime / availability
- Scalability (slow with large datasets)

👉 This is **how fast and stable the product is**.`

	case TopicUsabilityProductivity:
		return `Usability & Productivity

This is often confused with UX but different.

Topics:
- Workflow efficiency ("Takes 10 clicks to do X")
- Automation needs
- Keyboard shortcuts / power user features
- Defaults and smart suggestions
- Cognitive load ("Too many options, overwhelming")

👉 This is **how well the product helps users get work done**.`

	case TopicSecurityPrivacy:
		return `Security & Privacy

Very important in modern software.

Topics:
- Data leaks or privacy concerns
- Authentication / authorization issues
- Permission models (too permissive / too restrictive)
- Encryption and compliance (GDPR, HIPAA, etc.)
- Token handling

👉 This is **how safe the product is**.`

	case TopicCompatibilityIntegration:
		return `Compatibility & Integration

Big category for dev tools and enterprise software.

Topics:
- OS compatibility (Linux/Mac/Windows)
- Browser compatibility
- API integrations (GitHub, Slack, Stripe, etc.)
- Import/export formats
- Backward compatibility (breaking changes)
- SDK quality

👉 This is **how well it works with other systems**.`

	case TopicDeveloperExperience:
		return `Developer Experience (DX) *(especially for technical products)*

This is huge for tools like IDEs, SDKs, APIs.

Topics:
- Documentation quality
- API ergonomics
- Error messages clarity
- Debugging tools
- SDK stability
- Code examples and tutorials
- Language support

👉 This is **how pleasant it is to build on the product**.`

	case TopicPricingLicensing:
		return `Pricing & Licensing

Often overlooked in engineering feedback, but critical.

Topics:
- Too expensive
- Confusing pricing tiers
- Hidden costs
- Licensing restrictions
- Value-for-money perception
- Free tier limitations

👉 This is **how fair and transparent monetization feels**.`

	case TopicCustomerSupportCommunity:
		return `Customer Support & Community

Non-technical but important.

Topics:
- Support response time
- Quality of support answers
- Community activity (forums, Discord, GitHub issues)
- Bug triage speed
- Transparency in roadmap

👉 This is **how well the company supports users**.`

	case TopicInstallationSetupDeployment:
		return `Installation, Setup & Deployment

Critical for dev products and enterprise software.

Topics:
- Hard to install
- Broken installers
- Dependency hell
- Docker images missing
- CI/CD integration problems
- Upgrade/migration issues

👉 This is **how easy it is to get started and maintain**.`

	case TopicDataAnalyticsReporting:
		return `Data & Analytics / Reporting

For SaaS and enterprise apps.

Topics:
- Incorrect metrics
- Missing dashboards
- Hard-to-export reports
- Real-time vs delayed data
- Visualization issues

👉 This is **how well users can observe and analyze**.`

	case TopicLocalizationInternationalization:
		return `Localization & Internationalization

Often underrated.

Topics:
- Poor translations
- Missing languages
- Date/time/number formatting issues
- RTL layout bugs
- Cultural UX mismatches`

	case TopicProductStrategyRoadmap:
		return `Product Strategy & Roadmap Feedback

This is more high-level but real.

Topics:
- Direction dissatisfaction ("Why are you focusing on AI instead of stability?")
- Feature prioritization complaints
- Requests for enterprise vs consumer features
- Backward compatibility decisions`

	default:
		return ""
	}
}

// DisplayName returns a human-readable name for the topic.
func (t Topic) DisplayName() string {
	switch t {
	case TopicProductFunctionalityFeatures:
		return "Product Functionality & Features"
	case TopicUIUX:
		return "UI / UX"
	case TopicPerformanceReliability:
		return "Performance & Reliability"
	case TopicUsabilityProductivity:
		return "Usability & Productivity"
	case TopicSecurityPrivacy:
		return "Security & Privacy"
	case TopicCompatibilityIntegration:
		return "Compatibility & Integration"
	case TopicDeveloperExperience:
		return "Developer Experience"
	case TopicPricingLicensing:
		return "Pricing & Licensing"
	case TopicCustomerSupportCommunity:
		return "Customer Support & Community"
	case TopicInstallationSetupDeployment:
		return "Installation, Setup & Deployment"
	case TopicDataAnalyticsReporting:
		return "Data & Analytics / Reporting"
	case TopicLocalizationInternationalization:
		return "Localization & Internationalization"
	case TopicProductStrategyRoadmap:
		return "Product Strategy & Roadmap"
	default:
		return string(t)
	}
}

// AllTopics returns all available topics.
func AllTopics() []Topic {
	return []Topic{
		TopicProductFunctionalityFeatures,
		TopicUIUX,
		TopicPerformanceReliability,
		TopicUsabilityProductivity,
		TopicSecurityPrivacy,
		TopicCompatibilityIntegration,
		TopicDeveloperExperience,
		TopicPricingLicensing,
		TopicCustomerSupportCommunity,
		TopicInstallationSetupDeployment,
		TopicDataAnalyticsReporting,
		TopicLocalizationInternationalization,
		TopicProductStrategyRoadmap,
	}
}

// IsValid checks if the topic value is valid.
func (t Topic) IsValid() bool {
	switch t {
	case TopicProductFunctionalityFeatures,
		TopicUIUX,
		TopicPerformanceReliability,
		TopicUsabilityProductivity,
		TopicSecurityPrivacy,
		TopicCompatibilityIntegration,
		TopicDeveloperExperience,
		TopicPricingLicensing,
		TopicCustomerSupportCommunity,
		TopicInstallationSetupDeployment,
		TopicDataAnalyticsReporting,
		TopicLocalizationInternationalization,
		TopicProductStrategyRoadmap:
		return true
	default:
		return false
	}
}
