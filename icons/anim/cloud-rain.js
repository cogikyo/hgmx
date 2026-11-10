import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M4 14.899"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const drops = el.querySelectorAll('path[d$="v6"]');
  animate(drops, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
