import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll("rect");
  animate(bars, { scaleX: [0, 1.1, 1] }, { duration: τmid });

  const seams = el.querySelectorAll('path[d="M11 13v4"], path[d="M15 5v4"]');
  animate(seams, { scaleY: [0, 1] }, { duration: τmin, delay: τmid });
}
