import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const diaper = el.querySelectorAll('path[d^="M9 20"], path[d^="M2 13"], path[d^="M15 20"]');
  animate(diaper, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tabs = el.querySelectorAll('path[d="M2 9h4"], path[d="M22 9h-4"]');
  animate(tabs, { scaleX: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
