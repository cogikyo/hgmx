import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const olive = el.querySelector("ellipse");
  animate(olive, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pimento = el.querySelector("g");
  animate(pimento, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
