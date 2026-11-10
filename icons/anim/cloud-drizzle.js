import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M4 14.899"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const drops = ["M8 14", "M16 14", "M12 16", "M8 19", "M16 19", "M12 21"].map((d) => el.querySelector(`path[d^="${d}v"]`));
  animate(drops, { translateY: [-2, 0.5, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
