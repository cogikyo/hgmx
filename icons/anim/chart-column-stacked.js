import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll("rect");
  animate(bars, { scaleY: [0, 1.1, 1] }, { duration: τmid });

  const seams = el.querySelectorAll('path[d="M11 13H7"], path[d="M19 9h-4"]');
  animate(seams, { scaleX: [0, 1] }, { duration: τmin, delay: τmid });
}
