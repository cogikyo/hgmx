import { animate, τmax } from "../motion.js";

export default function (el) {
  const plug = el.querySelectorAll('path[d^="m2 22"], path[d^="M6.3"], path[d^="M7.5"], path[d^="M10.5"]');
  animate(plug, { translateX: [0, -2.5, 0.5, 0], translateY: [0, 2.5, -0.5, 0] }, { duration: τmax });

  const socket = el.querySelectorAll('path[d^="m19 5"], path[d^="m12 6"]');
  animate(socket, { translateX: [0, 2.5, -0.5, 0], translateY: [0, -2.5, 0.5, 0] }, { duration: τmax });
}
