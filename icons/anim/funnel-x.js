import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const funnel = el.querySelector('path[d^="M12.531"]');
  animate(funnel, { scaleX: [1, 0.8, 1.1, 1], scaleY: [1, 1.1, 0.95, 1] }, { duration: τmax });

  const bars = el.querySelectorAll('path[d^="m16.5"], path[d^="m21.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
