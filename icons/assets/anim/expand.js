import { animate, τmax } from "../motion.js";

export default function (el) {
  const [se, ne, sw, nw] = el.querySelectorAll("g");
  animate(se, { translateX: [0, 3, -1, 0], translateY: [0, 3, -1, 0] }, { duration: τmax });
  animate(ne, { translateX: [0, 3, -1, 0], translateY: [0, -3, 1, 0] }, { duration: τmax });
  animate(sw, { translateX: [0, -3, 1, 0], translateY: [0, 3, -1, 0] }, { duration: τmax });
  animate(nw, { translateX: [0, -3, 1, 0], translateY: [0, -3, 1, 0] }, { duration: τmax });
}
