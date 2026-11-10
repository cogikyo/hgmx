import { animate, τmax } from "../motion.js";

export default function (el) {
  const mustache = el.querySelector("path");
  animate(mustache, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });
}
