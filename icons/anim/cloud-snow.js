import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M4 14.899"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const flakes = ["M8 15", "M12 17", "M16 15", "M8 19", "M12 21", "M16 19"].map((d) => el.querySelector(`path[d^="${d}h"]`));
  animate(flakes, { translateY: [-3, 0], translateX: [1, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
