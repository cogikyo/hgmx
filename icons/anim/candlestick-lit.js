import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const candle = el.querySelectorAll('rect, path[d="m15 13-1-2"]');
  animate(candle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flame = el.querySelector('path[d^="M12 2s"]');
  animate(flame, { scaleY: [1, 0.8, 1.25, 1] }, { duration: τmid, delay: τmin });
}
