'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth-store';
import { AdminRoute } from '@/components/AdminRoute';
import { Feedback, Analysis, TopicStats } from '@/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { LogOut, Sparkles, Trash2, History } from 'lucide-react';

interface Statistics {
  total: number;
  averageRating: number;
  ratingDistribution: { [key: number]: number };
}

function AdminDashboard() {
  const router = useRouter();
  const { clearAuth } = useAuthStore();
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([]);
  const [statistics, setStatistics] = useState<Statistics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [latestAnalysis, setLatestAnalysis] = useState<Analysis | null>(null);
  const [isLoadingAnalysis, setIsLoadingAnalysis] = useState(false);
  const [topics, setTopics] = useState<TopicStats[]>([]);
  const [isLoadingTopics, setIsLoadingTopics] = useState(false);

  useEffect(() => {
    loadFeedbacks();
    loadLatestAnalysis();
    loadTopics();
  }, []);

  const loadLatestAnalysis = async () => {
    setIsLoadingAnalysis(true);
    try {
      const response = await apiClient.getLatestAnalysis();
      // Response might be wrapped in a data property or be direct
      const analysisData = response?.data || response;
      if (analysisData) {
        setLatestAnalysis(analysisData);
      }
    } catch (err: any) {
      if (err.response?.status !== 204) {
        // 204 means no analysis found, which is fine
        console.error('Failed to load latest analysis:', err);
      }
    } finally {
      setIsLoadingAnalysis(false);
    }
  };

  const loadTopics = async () => {
    setIsLoadingTopics(true);
    try {
      const response = await apiClient.getTopicsWithStats();
      const topicsData = response?.topics || response?.data?.topics || [];
      setTopics(topicsData);
    } catch (err: any) {
      console.error('Failed to load topics:', err);
    } finally {
      setIsLoadingTopics(false);
    }
  };

  const loadFeedbacks = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.listFeedbacks(1000, 0);
      const feedbacksData = response.feedbacks || response.data?.feedbacks || [];
      setFeedbacks(feedbacksData);

      // Calculate statistics
      const total = feedbacksData.length;
      const sum = feedbacksData.reduce((acc: number, f: Feedback) => acc + f.rating, 0);
      const averageRating = total > 0 ? sum / total : 0;

      const ratingDistribution: { [key: number]: number } = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 };
      feedbacksData.forEach((f: Feedback) => {
        ratingDistribution[f.rating] = (ratingDistribution[f.rating] || 0) + 1;
      });

      setStatistics({
        total,
        averageRating,
        ratingDistribution,
      });
    } catch (err: any) {
      if (err.response?.status === 401 || err.response?.status === 403) {
        clearAuth();
        router.push('/login');
      } else {
        setError('Failed to load feedbacks. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };


  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this feedback?')) {
      return;
    }

    try {
      await apiClient.deleteFeedback(id);
      await loadFeedbacks();
    } catch (err: any) {
      if (err.response?.status === 403) {
        setError('You do not have permission to delete feedbacks.');
      } else {
        setError('Failed to delete feedback. Please try again.');
      }
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

  return (
    <AdminRoute>
      <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="flex justify-between items-center">
            <h1 className="text-3xl font-bold">Admin Dashboard</h1>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => router.push('/feedback')}>
                Submit Feedback
              </Button>
              <Button variant="outline" onClick={() => router.push('/admin/analyses')}>
                <History className="h-4 w-4 mr-2" />
                Analysis History
              </Button>
              <Button variant="outline" onClick={() => { clearAuth(); router.push('/login'); }}>
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>

          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {/* Statistics */}
          {statistics && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Total Feedbacks</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-3xl font-bold">{statistics.total}</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Average Rating</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-3xl font-bold">{statistics.averageRating.toFixed(2)}/5.00</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Rating Distribution</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {[5, 4, 3, 2, 1].map((rating) => (
                      <div key={rating} className="flex items-center gap-2">
                        <span className="text-sm w-4">{rating}</span>
                        <div className="flex-1 bg-secondary rounded-full h-2">
                          <div
                            className="bg-primary h-2 rounded-full"
                            style={{
                              width: `${
                                statistics.total > 0
                                  ? (statistics.ratingDistribution[rating] / statistics.total) * 100
                                  : 0
                              }%`,
                            }}
                          />
                        </div>
                        <span className="text-sm w-8 text-right">
                          {statistics.ratingDistribution[rating]}
                        </span>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}

          {/* Latest AI Analysis */}
          <Card>
            <CardHeader>
              <div className="flex justify-between items-center">
                <div>
                  <CardTitle>Latest AI Analysis</CardTitle>
                  <CardDescription>Most recent analysis of feedback patterns and topics</CardDescription>
                </div>
                <Button variant="outline" onClick={() => router.push('/admin/analyses')}>
                  <History className="h-4 w-4 mr-2" />
                  View History
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {isLoadingAnalysis && (
                <p className="text-muted-foreground text-sm">Loading analysis...</p>
              )}
              {!isLoadingAnalysis && latestAnalysis && (
                <div className="space-y-6">
                  <div>
                    <h3 className="text-lg font-semibold mb-2">Summary</h3>
                    <div className="bg-muted rounded-md p-4 whitespace-pre-line text-sm">
                      {latestAnalysis.overall_summary}
                    </div>
                  </div>
                  <div className="flex gap-4">
                    <div>
                      <span className="text-sm text-muted-foreground">Sentiment: </span>
                      <Badge variant={latestAnalysis.sentiment === 'positive' ? 'default' : latestAnalysis.sentiment === 'negative' ? 'destructive' : 'secondary'}>
                        {latestAnalysis.sentiment}
                      </Badge>
                    </div>
                    <div>
                      <span className="text-sm text-muted-foreground">Status: </span>
                      <Badge variant={latestAnalysis.status === 'success' ? 'default' : latestAnalysis.status === 'failed' ? 'destructive' : 'secondary'}>
                        {latestAnalysis.status}
                      </Badge>
                    </div>
                    <div>
                      <span className="text-sm text-muted-foreground">Feedbacks Analyzed: </span>
                      <span className="font-semibold">{latestAnalysis.feedback_count}</span>
                    </div>
                  </div>
                  {latestAnalysis.key_insights && latestAnalysis.key_insights.length > 0 && (
                    <div>
                      <h3 className="text-lg font-semibold mb-2">Key Insights</h3>
                      <ul className="list-disc list-inside space-y-1 text-sm">
                        {latestAnalysis.key_insights.map((insight, i) => (
                          <li key={i}>{insight}</li>
                        ))}
                      </ul>
                    </div>
                  )}
                  <div className="pt-4 border-t">
                    <Button onClick={() => router.push(`/admin/analyses/${latestAnalysis.id}`)}>
                      View Full Analysis
                    </Button>
                  </div>
                </div>
              )}
              {!isLoadingAnalysis && !latestAnalysis && (
                <p className="text-muted-foreground text-sm">No analysis available yet. Analyses are generated automatically as feedbacks are submitted.</p>
              )}
            </CardContent>
          </Card>

          {/* Topics Section */}
          <Card>
            <CardHeader>
              <CardTitle>Topics</CardTitle>
              <CardDescription>Explore feedback by topic category</CardDescription>
            </CardHeader>
            <CardContent>
              {isLoadingTopics && (
                <p className="text-muted-foreground text-sm">Loading topics...</p>
              )}
              {!isLoadingTopics && topics.length > 0 && (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {topics.map((topic) => (
                    <Button
                      key={topic.topic}
                      variant="outline"
                      className="h-auto p-4 flex flex-col items-start justify-start hover:bg-accent transition-colors"
                      onClick={() => router.push(`/admin/topics/${topic.topic}`)}
                    >
                      <div className="w-full">
                        <h3 className="font-semibold text-left mb-2">{topic.topic_name}</h3>
                        <div className="flex items-center justify-between text-sm text-muted-foreground">
                          <span>{topic.feedback_count} feedback{topic.feedback_count !== 1 ? 's' : ''}</span>
                          {topic.average_rating > 0 && (
                            <span className="font-medium">
                              {topic.average_rating.toFixed(1)} ⭐
                            </span>
                          )}
                        </div>
                      </div>
                    </Button>
                  ))}
                </div>
              )}
              {!isLoadingTopics && topics.length === 0 && (
                <p className="text-muted-foreground text-sm">No topics available yet.</p>
              )}
            </CardContent>
          </Card>

          {/* Recent Submissions */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Submissions</CardTitle>
            </CardHeader>
            <CardContent>
              {feedbacks.length === 0 ? (
                <p className="text-muted-foreground">No feedbacks submitted yet.</p>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Rating</TableHead>
                      <TableHead>Comment</TableHead>
                      <TableHead>Date</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {feedbacks.slice(0, 20).map((feedback) => (
                      <TableRow key={feedback.id}>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <span className="text-lg">{'⭐'.repeat(feedback.rating)}</span>
                            <Badge variant="outline">{feedback.rating}</Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="max-w-md truncate">{feedback.comment}</div>
                        </TableCell>
                        <TableCell>
                          {new Date(feedback.created_at).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(feedback.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </AdminRoute>
  );
}

export default function AdminPage() {
  return <AdminDashboard />;
}
