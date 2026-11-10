import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cloud = el.querySelector('path[d^="M11 13"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0, 1], translateX: [-1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
