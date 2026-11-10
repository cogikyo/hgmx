import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cover = el.querySelectorAll('rect, path[d="M15 2v20"]');
  animate(cover, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tabs = el.querySelectorAll('path[d$="h5"]');
  animate(tabs, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
