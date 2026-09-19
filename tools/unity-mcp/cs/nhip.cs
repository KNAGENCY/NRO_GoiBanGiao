// Game co dang chay khung hinh khong? (Play ma Editor bi Pause, hoac Run In Background tat khi Unity mat focus -> game dung hinh)
var sb = new System.Text.StringBuilder();
sb.Append("isPlaying=").Append(UnityEditor.EditorApplication.isPlaying);
sb.Append(" isPaused=").Append(UnityEditor.EditorApplication.isPaused);
sb.Append(" frameCount=").Append(UnityEngine.Time.frameCount);
sb.Append(" time=").Append(UnityEngine.Time.time.ToString("F1"));
sb.Append(" runInBackground=").Append(UnityEngine.Application.runInBackground);
sb.Append(" focused=").Append(UnityEngine.Application.isFocused);
try { sb.Append(" gameTick=").Append(Game1.GameCanvas.gameTick); } catch (System.Exception e) { sb.Append(" gameTick EX ").Append(e.GetType().Name); }
try { sb.Append(" mainInstance=").Append(UnityEngine.Object.FindObjectOfType<Game1.Main>() != null); } catch (System.Exception e) { sb.Append(" main EX ").Append(e.GetType().Name); }
return sb.ToString();
