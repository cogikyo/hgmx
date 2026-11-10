import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const stem = el.querySelector('path[d="M12 8v12"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M6 4h12"]'), el.querySelector('path[d="M6 8h12"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
