import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M6 16.326"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  animate(cloud, { translateX: [0, -1, 1, 0] }, { duration: τmin, delay: τmid });

  const bolt = el.querySelector('path[d^="m13 12"]');
  animate(bolt, { pathLength: [0, 1], opacity: [0, 1, 0.25, 1] }, { duration: τmid, delay: τmin });
}
