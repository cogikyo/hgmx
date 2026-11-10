import { animate, τmax } from "../motion.js";

export default function (el) {
  const card = el.querySelector("g.card");
  animate(card, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmax });
}
