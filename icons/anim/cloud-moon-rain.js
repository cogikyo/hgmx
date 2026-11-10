import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M3 20"]');
  animate(cloud, { translateX: [-3, 1, 0], opacity: [0.5, 1] }, { duration: τmid });
  const moon = el.querySelector('path[d^="M18.376"]');
  animate(moon, { rotate: [0, -20, 8, 0], scale: [1, 0.85, 1] }, { duration: τmax });
  const drops = [el.querySelector('path[d^="M7 19"]'), el.querySelector('path[d^="M11 20"]')];
  animate(drops, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
