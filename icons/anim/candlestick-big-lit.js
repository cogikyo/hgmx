import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const candle = el.querySelectorAll('path[d^="M16 22"], path[d="M12 14v3"], path[d^="M17 17s"]');
  animate(candle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flame = el.querySelector('path[d^="M12 2S"]');
  animate(flame, { scaleY: [1, 0.8, 1.25, 1] }, { duration: τmid, delay: τmin });
}
