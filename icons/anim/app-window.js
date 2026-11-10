import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelectorAll('rect, path[d="M2 8h20"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tabs = [el.querySelector('path[d="M6 4v4"]'), el.querySelector('path[d="M10 4v4"]')];
  animate(tabs, { scaleY: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
