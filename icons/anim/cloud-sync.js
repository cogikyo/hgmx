import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M20.996"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const arrows = el.querySelector("g");
  animate(arrows, { rotate: [0, 180], scale: [1, 0.7, 1] }, { duration: τmid, delay: τmin });
}
