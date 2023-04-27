{
  outputs: [
    {
      name: 'Team',
      description: 'Teams in Product Development.',
      type_name: 'Custom["Team"]',
      source: {
        name: 'name',
        external_id: 'external_id',
      },
      attributes: [
        {
          id: 'description',
          name: 'Description',
          type: 'Text',
          source: 'description',
        },
        {
          id: 'goal',
          name: 'Goal',
          type: 'Text',
          source: 'goal',
        },
        {
          id: 'tech_lead',
          name: 'Tech lead',
          type: 'SlackUser',
          source: 'tech_lead',
        },
        {
          id: 'engineering_manager',
          name: 'Engineering manager',
          type: 'SlackUser',
          source: 'engineering_manager',
        },
        {
          id: 'product_manager',
          name: 'Product manager',
          type: 'SlackUser',
          source: 'product_manager',
        },
        {
          id: 'slack_user_group',
          name: 'Slack user group',
          type: 'SlackUserGroup',
          source: 'slack_user_group',
        },
        {
          id: 'slack_channel',
          name: 'Slack channel',
          type: 'String',
          source: 'slack_channel',
        },
        {
          id: 'alert_channel',
          name: 'Alert channel',
          type: 'String',
          source: 'alert_channel',
        },
        {
          id: 'linear_team',
          name: 'Linear team',
          type: 'LinearTeam',
          source: 'linear_team',
        },
        {
          id: 'members',
          name: 'Members',
          type: 'SlackUser',
          array: true,
          source: 'members',
        },
      ],
    },
  ],
  sources: [
    {
      inline: {
        entries: [
          {
            external_id: 'post-incident',
            name: 'Post-incident',
            description: 'Responsible for all product features that come after closing an incident.',
            goal: 'Increase learnings (decrease repeat incidents) for organisations.',
            tech_lead: 'alice',
            engineering_manager: 'bob',
            product_manager: 'carmen',
            slack_user_group: 'engineerz',
            slack_channel: '#team-post-incident',
            alert_channel: '#errors-pulse',
            linear_team: 'POS',
            members: [
              'dan',
              'ellie',
              'frank',
            ],
          },
          {
            external_id: 'response',
            name: 'Response',
            description: 'Responsible for all product features that power incident response.',
            goal: 'Make it easy to do the right thing (and quickly) during an incident.',
            tech_lead: 'george',
            engineering_manager: 'bob',
            product_manager: 'carmen',
            slack_user_group: 'engineerz',
            slack_channel: '#team-response',
            alert_channel: '#errors-pulse',
            linear_team: 'ENG',
            members: [
              'hannah',
              'isaac',
              'john',
            ],
          },
          {
            external_id: 'status-pages',
            name: 'Status pages',
            description: 'Responsible for the customer facing status pages product.',
            goal: 'Build the best status page product on the market (and annihilate the competition).',
            tech_lead: 'kyle',
            engineering_manager: 'liz',
            product_manager: 'carmen',
            slack_user_group: 'team-status-page',
            slack_channel: '#team-status-pages',
            alert_channel: '#errors-status-page-pulse',
            linear_team: 'SP',
            members: [
              'matthew',
              'norberto',
            ],
          },
        ],
      },
    },
  ],
}