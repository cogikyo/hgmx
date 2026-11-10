import { animate, τmax } from "../motion.js";

export default function (el) {
  const top = el.querySelector("g.top");
  animate(top, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const bottom = el.querySelector("g.bottom");
  animate(bottom, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmax });
}
