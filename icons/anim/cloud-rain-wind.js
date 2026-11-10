import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M4 14.899"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const drops = ["m9 13-", "m17 13-", "m9.2 22 "].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(drops, { translateX: [1.7, -0.4, 0], translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
