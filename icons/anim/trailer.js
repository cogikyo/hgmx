import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("g");
  const trailer = [body, el.querySelector("circle")];
  animate(trailer, { translateX: [-6, 2, 0], opacity: [0, 1] }, { duration: τmax });

  animate(body, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
