import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M13 16"]');
  animate(cloud, { translateX: [-3, 1, 0], opacity: [0.5, 1] }, { duration: τmid });
  const moon = el.querySelector('path[d^="M18.376"]');
  animate(moon, { rotate: [0, -20, 8, 0], scale: [1, 0.85, 1] }, { duration: τmax });
}
