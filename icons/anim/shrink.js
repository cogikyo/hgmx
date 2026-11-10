import { animate, τmax } from "../motion.js";

export default function (el) {
  const [se, sw, ne, nw] = el.querySelectorAll("path");
  animate(se, { translateX: [0, -1.5, 0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });
  animate(sw, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });
  animate(ne, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
  animate(nw, { translateX: [0, 1.5, -0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });
}
