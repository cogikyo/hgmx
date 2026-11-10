import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const [arrow, a] = el.querySelectorAll("g");
  animate(arrow, { translateY: [0, 4, -1, 0] }, { duration: τmax });

  const glyphs = [el.querySelector(":scope > path"), a];
  animate(glyphs, { opacity: [0, 1], y: [-3, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
