import { animate, τmax } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector("path");
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1], translateX: [-3, 1, 0] }, { duration: τmax });
}
