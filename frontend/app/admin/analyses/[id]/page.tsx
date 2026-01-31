'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth-store';
import { AdminRoute } from '@/components/AdminRoute';
import { AnalysisDetail, TopicAnalysis, FeedbackWithTopics } from '@/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { ArrowLeft, Calendar, TrendingUp, MessageSquare } from 'lucide-react';

const TOPIC_DISPLAY_NAMES: { [key: string]: string } = {
  product_functionality_features: 'Product Functionality & Features',
  ui_ux: 'UI / UX',
  performance_reliability: 'Performance & Reliability',
  usability_productivity: 'Usability & Productivity',
  security_privacy: 'Security & Privacy',
  compatibility_integration: 'Compatibility & Integration',
  developer_experience: 'Developer Experience',
  pricing_licensing: 'Pricing & Licensing',
  customer_support_community: 'Customer Support & Community',
  installation_setup_deployment: 'Installation, Setup & Deployment',
  data_analytics_reporting: 'Data & Analytics / Reporting',
  localization_internationalization: 'Localization & Internationalization',
  product_strategy_roadmap: 'Product Strategy & Roadmap',
};

function AnalysisDetailPage() {
  const router = useRouter();
  const params = useParams();
  const { clearAuth } = useAuthStore();
  const [analysisDetail, setAnalysisDetail] = useState<AnalysisDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const analysisId = params?.id as string;

  useEffect(() => {
    if (analysisId) {
      loadAnalysis();
    }
  }, [analysisId]);

  const loadAnalysis = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.getAnalysis(analysisId);
      // Response might be wrapped in a data property or be direct
      const responseData = response?.data || response;
      setAnalysisDetail(responseData);
    } catch (err: any) {
      if (err.response?.status === 401 || err.response?.status === 403) {
        clearAuth();
        router.push('/login');
      } else if (err.response?.status === 404) {
        setError('Analysis not found.');
      } else {
        setError('Failed to load analysis. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const getSentimentBadgeVariant = (sentiment: string) => {
    switch (sentiment) {
      case 'positive':
        return 'default';
      case 'negative':
        return 'destructive';
      case 'mixed':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'success':
        return 'default';
      case 'failed':
        return 'destructive';
      case 'processing':
        return 'secondary';
      default:
        return 'secondary';
    }
  };

  if (isLoading) {
    return (
      <AdminRoute>
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <h1 className="text-2xl font-bold mb-4">Loading...</h1>
          </div>
        </div>
      </AdminRoute>
    );
  }

  if (error || !analysisDetail) {
    return (
      <AdminRoute>
        <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
          <div className="max-w-7xl mx-auto">
            <Button variant="outline" onClick={() => router.push('/admin/analyses')} className="mb-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to History
            </Button>
            <Alert variant="destructive">
              <AlertDescription>{error || 'Analysis not found'}</AlertDescription>
            </Alert>
          </div>
        </div>
      </AdminRoute>
    );
  }

  const { analysis, topics, feedbacks } = analysisDetail;

  return (
    <AdminRoute>
      <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="flex items-center gap-4">
            <Button variant="outline" onClick={() => router.push('/admin/analyses')}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to History
            </Button>
            <h1 className="text-3xl font-bold">Analysis Details</h1>
          </div>

          {/* Analysis Overview */}
          <Card>
            <CardHeader>
              <div className="flex justify-between items-start">
                <div>
                  <CardTitle>Analysis Overview</CardTitle>
                  <CardDescription className="flex items-center gap-4 mt-2">
                    <span className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      Period: {new Date(analysis.period_start).toLocaleDateString()} -{' '}
                      {new Date(analysis.period_end).toLocaleDateString()}
                    </span>
                    {analysis.completed_at && (
                      <span className="flex items-center gap-1">
                        <TrendingUp className="h-4 w-4" />
                        Completed: {new Date(analysis.completed_at).toLocaleString()}
                      </span>
                    )}
                  </CardDescription>
                </div>
                <div className="flex gap-2">
                  <Badge variant={getSentimentBadgeVariant(analysis.sentiment)}>
                    {analysis.sentiment}
                  </Badge>
                  <Badge variant={getStatusBadgeVariant(analysis.status)}>
                    {analysis.status}
                  </Badge>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground mb-1">Overall Summary</p>
                  <p className="text-sm whitespace-pre-line">{analysis.overall_summary}</p>
                </div>
                {analysis.key_insights && analysis.key_insights.length > 0 && (
                  <div>
                    <p className="text-sm text-muted-foreground mb-2">Key Insights</p>
                    <ul className="list-disc list-inside space-y-1 text-sm">
                      {analysis.key_insights.map((insight, i) => (
                        <li key={i}>{insight}</li>
                      ))}
                    </ul>
                  </div>
                )}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t">
                  <div>
                    <p className="text-sm text-muted-foreground">Feedbacks Analyzed</p>
                    <p className="text-2xl font-bold">{analysis.feedback_count}</p>
                  </div>
                  {analysis.new_feedback_count !== null && analysis.new_feedback_count !== undefined && (
                    <div>
                      <p className="text-sm text-muted-foreground">New Feedbacks</p>
                      <p className="text-2xl font-bold">{analysis.new_feedback_count}</p>
                    </div>
                  )}
                  <div>
                    <p className="text-sm text-muted-foreground">Model</p>
                    <p className="text-lg font-semibold">{analysis.model}</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Tokens Used</p>
                    <p className="text-lg font-semibold">{analysis.tokens.toLocaleString()}</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Topics */}
          {topics && topics.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Topics Identified</CardTitle>
                <CardDescription>{topics.length} topics found in this analysis</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {topics.map((topic) => (
                    <Card key={topic.id}>
                      <CardHeader>
                        <div className="flex justify-between items-start">
                          <div>
                            <CardTitle className="text-base">{topic.topic_name}</CardTitle>
                            <CardDescription className="mt-1">
                              {topic.feedback_count} feedback{topic.feedback_count !== 1 ? 's' : ''}
                            </CardDescription>
                          </div>
                          <Badge variant={getSentimentBadgeVariant(topic.sentiment)}>
                            {topic.sentiment}
                          </Badge>
                        </div>
                      </CardHeader>
                      <CardContent>
                        <p className="text-sm whitespace-pre-line">{topic.summary}</p>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Analyzed Feedbacks */}
          {feedbacks && feedbacks.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle>Analyzed Feedbacks</CardTitle>
                <CardDescription>
                  {feedbacks.length} feedback{feedbacks.length !== 1 ? 's' : ''} analyzed in this analysis
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Rating</TableHead>
                      <TableHead>Comment</TableHead>
                      <TableHead>Topics</TableHead>
                      <TableHead>Date</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {feedbacks.map((feedback) => (
                      <TableRow key={feedback.id}>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <span className="text-lg">{'⭐'.repeat(feedback.rating)}</span>
                            <Badge variant="outline">{feedback.rating}</Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="max-w-md">{feedback.comment}</div>
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-wrap gap-1">
                            {feedback.topics.map((topicEnum, i) => (
                              <Badge key={i} variant="secondary" className="text-xs">
                                {TOPIC_DISPLAY_NAMES[topicEnum] || topicEnum}
                              </Badge>
                            ))}
                          </div>
                        </TableCell>
                        <TableCell>
                          {new Date(feedback.created_at).toLocaleDateString()}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </AdminRoute>
  );
}

export default function AnalysisDetailPageWrapper() {
  return <AnalysisDetailPage />;
}
